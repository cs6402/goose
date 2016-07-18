// hub
package control

type Hub struct {
	// Registered connections.
	connections map[string]*Conn

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Conn

	// Unregister requests from connections.
	unregister chan *Conn
}

var hub = Hub{
	broadcast:   make(chan []byte),
	register:    make(chan *Conn),
	unregister:  make(chan *Conn),
	connections: make(map[string]*Conn),
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn.id] = conn
		case conn := <-h.unregister:
			if _, ok := h.connections[conn.id]; ok {
				delete(h.connections, conn.id)
				close(conn.send)
			}
		case message := <-h.broadcast:
			for _, conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(h.connections, conn.id)
				}
			}
		}
	}
}
