package event

import (
    "sync"
    "fmt"
)

var (
    mgr = Manager{listeners:make(map[string][]Listener), singleListener:make(map[string]Listener)}
    ErrorListenerNotExist = fmt.Errorf("listener not found")
)

type Manager struct {
    mutex sync.RWMutex
    listeners map[string] []Listener
    mutexSingle sync.RWMutex
    singleListener map[string] Listener
}
type Listener interface {
    OnMessage(message string) (result string, err error)
}

func GetEventManager() *Manager {
    return &mgr
}

func (e *Manager) DispatchEvents(name string, message string) (result []string, err error) {
    e.mutex.RLock()
    defer e.mutex.RUnlock()
    if listeners, ok := e.listeners[name]; ok {
        for _, handler := range listeners {
            rs, c := handler.OnMessage(message)
            result = append(result, rs)
            err = c
            if err != nil {
                return
            }
        }
    } else {
        err = ErrorListenerNotExist
        return
    }
    return
}

func (e*Manager) AddListeners(name string, listener Listener) {
    e.mutex.Lock()
    defer e.mutex.Unlock()
    if listeners, ok := e.listeners[name]; ok {
        listeners = append(listeners, listener)
        e.listeners[name] = listeners
        return
    }
    var listeners []Listener
    listeners = append(listeners, listener)
    e.listeners[name] = listeners
}

func (e*Manager) RemoveListeners(name string, listener Listener) {
    e.mutex.Lock()
    defer e.mutex.Unlock()
    if listeners, ok := e.listeners[name]; ok {
        for i := 0; i < len(listeners); i++ {
            if listener == listeners[i] {
                listeners = append(listeners[:i], listeners[i+1:]...)
                e.listeners[name] = listeners
                return
            }
        }
    }
}

func (e *Manager) DispatchSingleEvent(name string, message string) (result string, err error) {
   e.mutexSingle.RLock()
   defer e.mutexSingle.RUnlock()
    if listener, ok := e.singleListener[name]; ok {
        return listener.OnMessage(message)
    } else {
        err = ErrorListenerNotExist
        return
    }
    return
}

func (e*Manager) AddSingleListener(name string, listener Listener) {
    e.mutexSingle.Lock()
    defer e.mutexSingle.Unlock()
    e.singleListener[name] = listener
}

func (e*Manager) RemoveSingleListener(name string) {
    e.mutexSingle.Lock()
    defer e.mutexSingle.Unlock()
    delete( e.singleListener, name)
}
