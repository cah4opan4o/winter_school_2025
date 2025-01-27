**Задание 1. Обзор консенсуса. Краткие вопросы**
1. В **папке** проекта создайте файл consensus_doc.go (или md).
2. Напишите **короткое описание** (5–10 предложений) “зачем нужен консенсус, чем отличается от лидер-элекции”.
3. Перечислите 3–4 ключевых свойства (Agreement, Validity, Liveness, No Double Value).
4. Укажите, что Raft / Paxos решают задачу при f < N/2 отказах.

**Задание 2. Спланировать структуру “mini-Raft”**
1. Создайте файл raft.go, где:
    * Тип RaftNode c полями:
```go
type RaftNode struct {
    id         int
    term       int
    state      string // "Follower", "Candidate", "Leader"
    log        []string
    votedFor   int
    leaderID   int
    // "commitIndex" "lastApplied" ...
    // каналы "inbox", "outbox"...
}
```
2. Решите, **как** будете пересылать RPC:
    * RequestVote (candidate -> всех)
    * AppendEntries (leader -> followers)
3. Заведите type Message struct { kind string; term int; from, to int; entry string; ... }.

**Задание 3. Добавить таймеры (follower timeout, heartbeat)**
1. В Raft узел, будучи follower, **ждёт** heartbeat от лидера. Если нет heartbeat в течение “ElectionTimeout” (random 150–300 ms в реале, мы эмулируем).
2. Используем “горутина + time.AfterFunc(…)” (или time.Sleep(...)).
3. При срабатывании таймера: follower → “Candidate”, увеличивает term, RequestVote.
    * Продумайте, как узел “блокируется” в select {} c каналом “heartbeatChan” и/или “time.After” ~ 2s.
    * Если “heartbeatChan” сработал, узел “обновляет” всё. Если “time.After” сработал раньше — переходим в candidate.

**Задание 4. Расширить структуру “RaftNode” для лога и AppendEntries**
1. В созданной структуре RaftNode c log []string, замените log []string на log []LogEntry:
```go
type LogEntry struct {
    Term  int
    Command string
}
```
2. Определите новый тип сообщения:
```go
type AppendEntriesArgs struct {
    Term     int
    LeaderID int
    PrevLogIndex int
    PrevLogTerm  int
    Entries []LogEntry
    LeaderCommit int
}  

type AppendEntriesReply struct {
    Term       int
    Success    bool
    MatchIndex int
}
```
3. Для упрощения, можно игнорировать prevLogIndex/Term и **предположить**, что follower логи идентичны. Всё равно создадим структуру, но часть полей можем не использовать.

**Задание 5. Реализовать логику “AppendEntries” (leader → followers)**
1. **Лидер**, получив новую команду "commandX" (условно), делает:
```go
rn.log = append(rn.log, LogEntry{Term: rn.currentTerm, Command: "commandX"})
// затем рассылает AppendEntries
for each follower in cluster {
    send AppendEntriesArgs{Term: rn.currentTerm, LeaderID: rn.id, Entries: []LogEntry{{Term: rn.currentTerm, Command: "commandX"}}, ...} to follower
}
```
2. **Follower** при получении AppendEntriesArgs:
    * Сравнивает Term с rn.currentTerm; если Term <, то Reply{Success:false}.
    * Иначе: rn.currentTerm = Term.
    * Append к своему log “Entries…”.
    * Возвращает AppendEntriesReply{Term: rn.currentTerm, Success:true, MatchIndex = lastIndexOfLog}.
3. **Leader** при получении AppendEntriesReply от **большинства** follower-ов “подтверждает” запись.
    * Для упрощения: можете хранить matchIndex\[followerID\]. Если “Success=true, matchIndex=…”, обновляйте. Когда “кол-во matchIndex >= majority”, можно считать “команда зафиксирована”.

**Задание 6. “commitIndex” и оповещение о применении**
1. Добавьте в RaftNode поле commitIndex int.
2. Когда лидер видит, что “большинство” имеют “matchIndex >= X”, он обновляет commitIndex = X.
3. Можем делать “apply” на локальное состояние (например, fmt.Println("Apply command:", rn.log\[X\].Command)) — учебно.
4. _Follower_, получая “LeaderCommit” в “AppendEntriesArgs.LeaderCommit”, если leaderCommit > rn.commitIndex, обновляет rn.commitIndex = leaderCommit, и “применяет” записи логов до этой позиции.

**Задание 7. Имитация простого “конфликта лога”**
1. **Смоделируйте** ситуацию:
    * Лидер (узел L) получил команду “A”, записал LogEntry(term=2, cmd="A" ), разослал **не всем** (скажем, 1 follower принял, другой — нет, т.к. “сообщение потерялось”).
    * Лидер “упал” (alive=false).
2. **Новый лидер**:
    * Другой узел (F) со “старым” логом (без “A”) становится лидером по RequestVote (т.к. большинство “видит” его).
    * Он начинает “AppendEntries” → follower (который видел “A”) теперь увидит “несовпадение” в логе.
3. **Реализуйте** (минимально) ответ follower: если prevLogIndex, prevLogTerm не совпадает с тем, что у него в логе, → “success=false”. Лидер тогда уменьшит prevLogIndex (или “truncateConflict”) и пошлёт заново.
4. **Результат**: запись “A” потеряется, т.к. новый лидер не знал о ней.
5. Посмотрите, что в итоге все узлы синхронизируются с новым лидером.

**Задание 8. Демонстрация результатов дня main_day4.go:**
1. Создайте 3–5 RaftNode:
    * Изначально state=“Follower”, term=0, leaderID=-1, votedFor=-1, log=[].
    * Запускаете горутины.
2. В горутине func (rn *RaftNode) run():
    * “for { select { case msg := <-rn.inbox: … } }”
    * Реагируйте на RequestVote (если msg.term >= rn.term, rn.votedFor = msg.from, …).
    * Реагируйте на AppendEntries (heartbeat).
3. Попробуйте **один** узел “не получать heartbeat”, и по истечению timeout → становится candidate, рассылает RequestVote…
    * Если получает большинство голосов, становится Leader → рассылает AppendEntries heartbeat.
4. В run() при каждом событии (RequestVote, AppendEntries, “timeout triggers candidate”) делайте fmt.Println(...).
5. Наблюдайте, что в итоге **один** узел становится Leader, другие — Follower, term поднимается.
6. Запускаем горутины, “ElectionTimeout” → узел становится candidate, выиграл → leader.
7. “Клиент” (в main()) отправляет commandX, commandY, ... лидеру.
    * Лидер рассылает AppendEntries, follower отвечает. Лидер печатает “команда зафиксирована!”
8. Покажите, что если лидер упал **до** того, как команда была подтверждена большинством, она не станет зафиксированной:
    * Лидер “сохраняет” запись в локальный лог.
    * До того, как придёт ACK от большинства, вызывайте leader.alive = false.
    * Новый лидер не узнает о записи и “лог перезапишется”.
    * **Выведите**: “Запись X потерялась, т.к. не была зафиксирована”.
9. Вместо прямой outbox\[to\] <- msg, сделайте “сеть” — goroutine, которая с вероятностью X% теряет пакет, или time.Sleep(...) для имитации задержки.
10. Проверьте, что Raft всё равно стабилизируется (если большинство узлов + лидер действительно видят друг друга).
11. **Отладочный вывод**: “Message from X to Y was dropped”.
