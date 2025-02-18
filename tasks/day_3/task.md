**Задание 1. Создать каркас «узла» (Node)**
1. Создайте файл leader_election.go (или bully.go).
2. Опишите структуру:
```go
type Node struct {
    id       int
    isLeader bool
    alive    bool
    // какие-то поля для каналов, или ссылки на другие узлы
}
```
3. Подумайте, как организовать «сеть сообщений». Варианты:
    * **Вариант A**: для каждого узла объявить канал входящих сообщений inbox chan Message, где Message — type struct { from, to, kind string; data any }.
    * **Вариант B**: сделать «глобальный диспетчер» (map[id]chan Message), где каждый узел «слушает» свой канал.
4. Инициализируйте несколько Node c разными ID, скажем, 3–5 узлов.

**Задание 2. Реализовать Bully-алгоритм (базовая логика)**
1. **Процедура «узел замечает отсутствие лидера»**:
    * Узел node.id = X рассылает Message{ from: X, kind: "ELECTION" } всем узлам, у кого id > X.
2. **Ответ «OK»**:
    * Узел с id > X, получив ELECTION от X, отсылает OK обратно.
    * Узел-инициатор, получив OK, понимает, что есть «кандидат» с ID больше, и «ждёт» (или прекращает выбор).
3. **Если не получил «OK»** (по таймеру?), узел объявляет себя лидером, рассылает COORDINATOR.
4. **COORDINATOR**:
    * Все узлы, получив COORDINATOR, записывают leaderID = <посланный>.

**Задание 3. Создать “кольцо” узлов**
1. Расширьте структуру:
```go
type Node struct {
    id       int
    nextID   int  // кому пересылать в кольце
    leaderID int
    alive    bool
    inbox    chan Message
    // ...
}
```
2. Инициализируйте N узлов c ID = 0..N-1; свяжите их в кольцо: node[i].nextID = (i+1) % N.

**Задание 4. Реализовать Ring-based election логику**
1. **Тип сообщения** Message может включать:
```go
type Message struct {
    kind   string // "ELECTION", "COORDINATOR"
    ids    []int  // список ID, через который проходит сообщение
    fromID int
}
```
* Или хранить единственный “maxID” вместо массива. В классическом Chang–Roberts алгоритме чаще собирают все ID, но нам достаточно запоминать **текущий максимум**.
2. **Запуск** ELECTION:
    * Если узел замечает, что leaderID = -1, а сам считает, что «надо выбрать лидера», он создаёт сообщение Message{kind:"ELECTION", ids: []int{id}}, отправляет node[nextID].inbox <- msg.
3. **Обработка** ELECTION:
    * Если msg.ids (или msg.maxID) меньше node.id, заменить на node.id.
    * Переслать дальше по кольцу: node[nextID].inbox <- msg.
    * Если узел видит, что len(msg.ids) != 0 и msg.ids[0] == node.id (сообщение вернулось к инициатору), значит maxID внутри msg — победитель. Узел рассылает COORDINATOR.
4. **COORDINATOR**:
    * Передаёт leaderID = maxID по кольцу, пока не дойдёт до инициатора снова; все записывают leaderID.

**Задание 5. «Подготовить» ведущий узел (лидера) к координации**
1. В коде из предыдущих частей (Bully и Ring-based) после определения лидера **каждый** узел знает leaderID.
2. Узлу с id == leaderID нужно добавить функциональность:
```go
func (n *Node) StartGlobalCollection() {
    // Рассылать запрос "COLLECT" другим узлам
    // И ждать ответов
}
```
3. Остальные узлы должны уметь обрабатывать Message{ kind: "COLLECT" }:
    * “COLLECT” может включать “request data”.
    * Узел отвечает Message{kind: "COLLECT_REPLY", data: localData} лидеру.

**Задание 6. «Локальные данные» каждого узла**
1. Допустим, у каждого узла своя “часть графа” (или просто «число пользователей = 100»).
2. Мы хотим сложить эти числа.
3. В структуре Node добавьте поле localCount int. Инициализируйте (например, rand.Intn(50) + 50, чтобы у каждого узла была “своя” цифра).
4. Когда узел получает COLLECT от лидера:
    * Формирует Message{ kind:"COLLECT_REPLY", fromID: node.id, data: node.localCount }.
    * Отправляет обратно leaderID.inbox <- ....

**Задание 7. Логика лидера: «Собрать ответы»**
1. Лидер, вызвав StartGlobalCollection(), рассылает Message{kind:"COLLECT", fromID: leaderID} всем узлам (кроме себя).
2. Лидер должен ожидать **N-1** ответов (если узлов N, сам себе не отправляет).
3. Заводим счётчик received = 0, sum = 0.
4. Когда приходит COLLECT_REPLY, лидер делает sum += msg.data.(int); received++.
5. Когда received == N-1, лидер объявляет итог sum.
6. (Упрощение) Игнорируем сбои узлов на данном этапе. Или (опционально) учитываем, что узел не отвечает, тогда по таймауту лидер решает, что “узел упал”.

**Задание 8. Демонстрация результатов третьего дня**
1. Создайте day3_main.go, где:
    * Инициализируете 3–5 узлов (ID = 1,2,3..).
    * Организуете каналы inboxChan[id] для каждого.
    * Запускаете горутину на узел, в которой for { select { case msg := <-inboxChan[id]: ... } }.
    * «Имитация»: пусть сразу никто не знает лидера. Узел 1 решает запустить ELECTION.
    * Сымитируйте отсутствие лидера (leaderID = -1), пусть узел 0 отправит ELECTION (или “заметит” отсутствие).
    * Посмотрите, как сообщение “гуляет” по кольцу, кто становится лидером (узел с **max ID**).
    * После выбора лидера (leaderID), вызываем у лидера StartGlobalCollection().
    * Печатаем итог: “Лидер X собрал сумму Y”.
2. Посмотрите, как приходят ответы, кто станет лидером. Выведите итог.
3. **Имитируйте сбой**: один узел не ответил (alive=false). Лидер ждёт (timeout).
4. **Объедините**: лидер, после сбоя, переизбирается. Новый лидер потом собирает данные.