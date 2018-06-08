//implements a generic stack based off linked-list design, then provides a layer specifically for a stack of matrices on top of that
package main
import (
    "sync"
)
type node struct {
    //cargo can be any datatype that implements any or no methods
    cargo interface{}
    next *node
}

//a stack is simply a pointer to the head plus a Mutex lock (this makes it so only one thread can access the stack at a time and is atomically updated much like a semaphore. Honestly, there's a 0 percent chance of thread-safety related issues, since creating 4x4 matrices takes an insignificant amount of time compared to multiplying it, plus the pop() and push() functions are run in a single thread anyways since order matters. It's fun to try out, anyways


type Stack struct {
    head *node
    lock sync.Mutex
}


func (stack *Stack) Push(data interface{}) {
    stack.lock.Lock()

    newHead := node{data, stack.head}
    stack.head = &newHead

    stack.lock.Unlock()
}

func (stack *Stack) Peek() interface{} {
    var data interface{} = nil
    stack.lock.Lock()
    head := stack.head
    if (head != nil) {
        data = head.cargo
    }
    stack.lock.Unlock()
    return data
}

func (stack *Stack) Pop() interface{} {
    var data interface{} = nil
    stack.lock.Lock()
    head := stack.head
    if (head != nil) {
        data = head.cargo
        stack.head = head.next
    }
    stack.lock.Unlock()
    return data
}



//matrix specific

func MakeWorldStack() *Stack {
    mat := IdentityMat(4)
    stack := Stack{&node{&mat,nil},sync.Mutex{}}
    return &stack
}

func (stack *Stack) GetWorld() *Matrix {
    var i interface{} = stack.Peek()
    mat := i.(*Matrix)
    return mat
}

func (stack *Stack) PopWorld() * Matrix {
    var i interface{} = stack.Pop()
    mat := i.(*Matrix)
    return mat
}

func (stack *Stack) PushWorld() {
    //make a copy
    worldCopy := (DeepCopy(stack.GetWorld()))
    stack.Push(&worldCopy)
}
