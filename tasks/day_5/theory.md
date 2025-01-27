1. Расширенное использование goroutines и каналов для MapReduce

1.1. Каналы «задач» и «результатов»
•	В простой учебной MapReduce без «container/queue», мы имеем:

type MapTask struct {
ChunkID int
Data    []byte // или string
}
type MapResult struct {
ChunkID int
Pairs   []Pair // (k,v)
}

	•	master может иметь mapTaskChan chan MapTask для “выдать задачу”,
	•	и mapResultChan chan MapResult для “получить результат”.

	•	Аналогично для reduce:

type ReduceTask struct {
Key     string
Values  []int
}
type ReduceResult struct {
Key     string
Result  int
}

	•	reduceTaskChan, reduceResultChan.

1.2. Per-Task goroutine vs Worker goroutine
•	Вариант A: “pool” воркеров – N горутин, каждая “ждёт” MapTask, обрабатывает → посылает MapResult.
•	Вариант B: на каждую задачу “go func() { … }()”. Но это может быть слишком много горутин, если задач очень много.

(Обычно вариант A, “Worker goroutines” + каналы.)

1.3. Учёт сбойных воркеров (re-run Task)
•	Master может задать “timeout”: если в течение X не пришёл MapResult, worker считаем “упавшим”. Назначаем задачу другому worker.
•	Всё реализуем вручную (slice или map для “taskStatus”), без готовых структур.

2. Combiner в MapReduce (локальное агрегирование)

2.1. “mapFunc” vs “combinerFunc”
•	При классическом WordCount:

func mapFunc(data string) []Pair {
// split на слова
// emit (word, 1)
}
func combinerFunc(pairs []Pair) []Pair {
// локально агрегируем (word -> sum) -> emit (word, sum)
}


	•	Можно “combinerFunc” вызвать прямо внутри mapWorker, до отправки “MapResult” → сократится shuffle.

2.2. Организация combiner
•	В worker, после mapFunc:

rawPairs := mapFunc(task.Data)
combined := localCombine(rawPairs) // создаём map[word] => sum, потом -> []Pair
// затем Master получает эти combined

3. Интеграция с “лидер-элекция” / “mini-Raft”

3.1. “Master = Leader”
•	Если у вас есть лидер (Bully / RaftLeader), он становится Master MapReduce:
•	Рассылает MapTasks → собирает результаты → Shuffle → ReduceTasks.
•	Если лидер упал, новый лидер может «перенять»? Нужно хранить состояние (какие таски завершены) – возможно, “mini-Raft log”.

3.2. Raft log “Map job” state
•	Продвинутый вариант: при смене лидера, новый лидер восстанавливает, какие chunk’и уже map’лены, какие reduce’ы сделаны.
•	Это требует “подписать” в Raft log: “mapTask i done”, “reduceTask j done”.

4. Расширенное тестирование MapReduce

4.1. “Стресс” тест
•	Сгенерировать “большую” строку (200–500 KB), разбить на 5–10 chunks, эмулируем Map (в задаче WordCount).
•	Запускаем 3–5 worker’ов.
•	Считаем, всё ли корректно суммируется.

4.2. Сбой worker
•	Как выше упоминали, master ждёт MapResult, если нет – re-assign.
•	Worker “alive = false”, “run()” остановилась.
•	Снова “mapTaskChan <- task” для другого worker.

4.3. “Потери” shuffle
•	Можно “shuffleChan” терять часть (k,v)? Master/Reduce re-requests?
•	В реальной Hadoop – много деталей, здесь учебно можно пропустить.

5. Профилирование и оптимизация

5.1. “pprof” (необязательно)
•	Можно import _ "net/http/pprof" и запускать go tool pprof.

5.2. Срезы vs map
•	При combiner: “map[string]int” локально, потом конвертировать в []Pair.

5.3. Sync primitives
•	Используем sync.Mutex, sync.WaitGroup.

6. Финальное связывание: “распределённый соцграф” + MapReduce
   •	Идея:
    1.	У нас есть N узлов, каждый хранит часть “соцграфа” (срез вершин, рёбер).
    2.	“Master” решает: “подсчитаем, сколько пользователей с возраст>30” (условно).
    3.	Рассылает mapTask (каждому узлу – “пробегись по своему chunk-у, найди (age>30,1)”).
    4.	Combiner локально суммирует, master собирает, shuffle (хотя ключ=“age>30” един).
    5.	“reduceTask” → общая сумма.
          •	No “mapreduce” pkg, всё ручками (goroutines, channels, slices).

Рекомендации / Напоминания
1.	Будьте аккуратны с каналами: при goroutine “worker” завершении, канал нужно закрыть или как-то обрабатывать.
2.	Учебная реализация – не нужно писать “production-level retry” для shuffle, но можно продемонстрировать 1-2 timeouts/сбоев.
3.	Код может стать громоздким – разбейте на файлы: mapreduce/master.go, mapreduce/worker.go, mapreduce/common.go.