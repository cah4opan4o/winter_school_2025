package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message — структура для передачи сообщений между узлами
type Message struct {
	kind   string // "ELECTION", "OK", "COORDINATOR", "COLLECT", "COLLECT_REPLY"
	fromID int
	data   any
}

// Node — структура узла в сети
type Node struct {
	id        int
	nextID    int
	leaderID  int
	alive     bool
	localData int
	inbox     chan Message
	nodes     map[int]*Node
	mu        sync.Mutex
}

func CreateNodes(n int) map[int]*Node {
	nodes := make(map[int]*Node)

	for i := 0; i < n; i++ {
		nodes[i] = &Node{
			id:        i,
			nextID:    (i + 1) % n,
			leaderID:  -1, // неизвестный лидер
			alive:     true,
			localData: rand.Intn(50) + 50, // случайные данные
			inbox:     make(chan Message, 10),
			nodes:     nodes,
		}
	}

	// Запускаем горутину для каждого узла
	for _, node := range nodes {
		go node.Listen()
	}

	return nodes
}

func (n *Node) StartElection() {
	fmt.Printf("Node %d: Начинаю выборы...\n", n.id)

	// Отправляем "ELECTION" всем с большим ID
	higherNodes := false
	for id := range n.nodes {
		if id > n.id && n.nodes[id].alive {
			n.nodes[id].inbox <- Message{kind: "ELECTION", fromID: n.id}
			higherNodes = true
		}
	}

	// Если не получили "OK", становимся лидером
	if !higherNodes {
		n.BecomeLeader()
	}
}

func (n *Node) BecomeLeader() {
	n.leaderID = n.id
	fmt.Printf("Node %d: Я лидер!\n", n.id)

	// Рассылаем "COORDINATOR"
	for _, node := range n.nodes {
		if node.id != n.id && node.alive {
			node.inbox <- Message{kind: "COORDINATOR", fromID: n.id}
		}
	}
}

func (n *Node) HandleElection(msg Message) {
	if msg.fromID < n.id {
		fmt.Printf("Node %d: Отправляю OK %d\n", n.id, msg.fromID)
		n.nodes[msg.fromID].inbox <- Message{kind: "OK", fromID: n.id}
		n.StartElection()
	}
}

func (n *Node) HandleCoordinator(msg Message) {
	n.leaderID = msg.fromID
	fmt.Printf("Node %d: Новый лидер %d\n", n.id, n.leaderID)
}

func (n *Node) Listen() {
	for {
		select {
		case msg := <-n.inbox:
			switch msg.kind {
			case "ELECTION":
				n.HandleElection(msg)
			case "OK":
				// Узел понял, что есть более приоритетный кандидат
			case "COORDINATOR":
				n.HandleCoordinator(msg)
			case "COLLECT":
				n.HandleCollect(msg)
			case "COLLECT_REPLY":
				n.HandleCollectReply(msg)
			}
		}
	}
}

func (n *Node) StartRingElection() {
	msg := Message{kind: "ELECTION", data: n.id, fromID: n.id}
	n.nodes[n.nextID].inbox <- msg
}

func (n *Node) HandleRingElection(msg Message) {
	maxID := msg.data.(int)

	if maxID < n.id {
		msg.data = n.id
	}

	if msg.fromID == n.id {
		// Сообщение вернулось, значит, мы лидер
		n.BecomeLeader()
	} else {
		n.nodes[n.nextID].inbox <- msg
	}
}

func (n *Node) StartGlobalCollection() {
	fmt.Printf("Node %d: Начинаю сбор данных...\n", n.id)
	for _, node := range n.nodes {
		if node.id != n.id {
			node.inbox <- Message{kind: "COLLECT", fromID: n.id}
		}
	}
}

func (n *Node) HandleCollect(msg Message) {
	n.nodes[msg.fromID].inbox <- Message{kind: "COLLECT_REPLY", fromID: n.id, data: n.localData}
}

func (n *Node) HandleCollectReply(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.localData += msg.data.(int)
	fmt.Printf("Node %d: Получил %d от Node %d. Итоговая сумма: %d\n", n.id, msg.data, msg.fromID, n.localData)
}

func main() {
	nodes := CreateNodes(5)

	// Имитация сбоя узла
	time.Sleep(2 * time.Second)
	nodes[4].alive = false
	nodes[2].StartElection()

	// Дожидаемся выбора лидера
	time.Sleep(3 * time.Second)

	// Запускаем сбор данных
	leader := nodes[nodes[2].leaderID]
	leader.StartGlobalCollection()

	// Даем время на завершение сбора
	time.Sleep(5 * time.Second)
}
