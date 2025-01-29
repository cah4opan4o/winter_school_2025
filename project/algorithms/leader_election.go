package algorithms

type Node struct{
	id int
    // nextID int  // кому пересылать в кольце
    // leaderID int
    isLeader bool
    alive bool
    inbox chan Message
}

func NewNode(id int)*Node{
	return &Node{id: id,isLeader: false, alive: true, inbox: chan Message}
}

type Message struct {
    kind string // "ELECTION", "COORDINATOR"
    ids []int  // список ID, через который проходит сообщение
    fromID int
}

func NewMessage(kind string, ids[] int ,fromID int)*Message{
	return &Message{kind: kind, ids: append(ids,ids) }
}


func (n* Node) Bully(x int,spisok []Node)(){
	
	return 
}