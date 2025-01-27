1. Работа с горутинами (Go concurrency basics)

1.1. Запуск горутины
•	Любая функция (или анонимная лямбда) может быть запущена как горутина:

go func() {
// код
}()


	•	Выполнение асинхронное: основная горутина не ждёт окончания, unless вы делаете sync.WaitGroup или каналы.

1.2. Каналы
•	Создание канала: ch := make(chan int) (не буферизированный) или chBuf := make(chan int, 10) (буфер = 10).
•	Отправка: ch <- x
•	Получение: val := <-ch
•	Закрытие: close(ch), после чего читатели могут получать val, ok := <-ch, где ok == false означает канал закрыт.

Пример “pinger–ponger”:

func pinger(ch chan<- string) {
for {
ch <- "ping"
}
}

func ponger(ch chan<- string) {
for {
ch <- "pong"
}
}

func printer(ch <-chan string) {
for msg := range ch {
fmt.Println(msg)
}
}

func main() {
c := make(chan string)
go pinger(c)
go ponger(c)
go printer(c)
time.Sleep(1 * time.Second)
}

	•	Мы не используем готовых контейнеров (list, heap), но каналы – часть языка.

1.3. Селект (select)
•	Позволяет обрабатывать сразу несколько операций чтения/записи в каналы, таймеры, т.д.:

select {
case msg := <-ch1:
fmt.Println("Got from ch1:", msg)
case ch2 <- "hello":
fmt.Println("Sent to ch2")
case <-time.After(2 * time.Second):
fmt.Println("Timeout!")
default:
fmt.Println("No activity")
}


	•	Полезно при реализации таймаутов (в лидер-элекции, mini-Raft).

2. Таймеры, time.After / time.Ticker

Для алгоритмов лидер-элекции (Bully, Ring) и mini-Raft вам понадобятся таймеры:
1.	time.After(d): возвращает канал, который “стрельнёт” через d:

select {
case <-time.After(3 * time.Second):
fmt.Println("Timeout, become candidate!")
case msg := <-inbox:
// handle message
}

	2.	time.NewTicker(d): канал, который “стреляет” каждые d:

ticker := time.NewTicker(1 * time.Second)
go func() {
for t := range ticker.C {
fmt.Println("Tick at", t)
}
}()

	•	Внимание: не забывать ticker.Stop() если не нужно больше тиков.

3. Организация “распределённых” узлов без стандартных контейнеров

3.1. Структура узла (Node)

В “Лидер-элекции” (Bully / Ring) или “mini-Raft” вы зачастую делаете:

type Node struct {
id         int
alive      bool
inbox      chan Message
leaderID   int
// Для Bully: no big container, just slice of other nodeIDs
// Для Ring: nextID int (кольцо)
// Для Raft: term, state, log ...
}

	•	Message struct (kind string, from, to int, …) позволяет различать типы сообщений ("ELECTION", "OK", "APPEND_ENTRIES", etc.).
	•	Срез []Node или map[int]*Node (где ключ = id) может хранить “кластер”.

3.2. “Транспорт”

Вместо готовых RPC-библиотек, вы вручную делаете что-то вроде:

func sendMessage(cluster map[int]*Node, msg Message) {
if node, ok := cluster[msg.to]; ok && node.alive {
node.inbox <- msg
}
}

	•	Или “сеть” goroutine, которая может имитировать потери.

3.3. “Учебная” обработка: select { case msg := <-n.inbox: … }

Внутри Node.run():

func (n *Node) run(cluster map[int]*Node) {
for n.alive {
select {
case msg := <-n.inbox:
// switch msg.kind ...
case <-time.After(2 * time.Second):
// election timeout
}
}
}

4. Логика лидер-элекции (Bully / Ring)

•	Bully: держим []int всех ID, если id < other рассылаем “ELECTION”, “OK”, “COORDINATOR” просто через каналы.
•	Ring: просто поле nextID int, в Node (указывающее, кому переслать).
•	Реализация: “протокол” обрабатывается switch-case в select.

5. “mini-Raft” (подробнее)

Тут ещё глубже нужны каналы, таймеры:

5.1. Структура

type RaftNode struct {
id       int
term     int
state    string // "Follower"/"Candidate"/"Leader"
log      []LogEntry
commitIndex int
// ...
inbox chan RaftMessage
}

	•	LogEntry – просто type LogEntry struct { Term int; Command string }.

5.2. “RPC” (псевдо)

type AppendEntriesArgs struct {
Term     int
LeaderID int
Entries  []LogEntry
// ...
}
type AppendEntriesReply struct {
Term    int
Success bool
}

	•	При “send” -> “receive” всё через inbox chan.
	•	“select / case <-n.inbox” → parse msg.kind.

5.3. Timer usage
•	“ElectionTimeout” → if no AppendEntries → become candidate.
•	time.AfterFunc(d, func(){ … }) or select{case <-time.After(d):...}.

7. Советы по отладке и профилированию
    1.	fmt.Println – основной инструмент логирования.
    2.	Вычислительные горутины: при большом количестве горутин, можно пересмотреть GOMAXPROCS, но обычно 1..4 логических CPU достаточно.
    3.	Race: если общий доступ к разделяемым структурам, можно запустить go run -race main.go. Но мы стараемся минимизировать общий доступ, используя каналы.

8. Итог: «дополнительные» Go возможности для 3-го дня
   •	Concurrency (горoutines, channels, select, time.After, time.Ticker) — ключевые для реализации распределённых алгоритмов (лидер-элекция, Raft).
   •	Struct + slice / map — всё ещё основа (никаких “container/list/heap”).
   •	No ready-made RPC libs: вручную делаем message-passing, “sendMessage()” / “inbox chan”.
   •	Timers для “timeout” (Bully, Ring, Raft).