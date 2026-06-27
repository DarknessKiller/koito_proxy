type ToastType = 'success' | 'error' | 'info'

interface Toast {
  id: number
  message: string
  type: ToastType
}

class ToastManager {
  items = $state<Toast[]>([])
  #nextId = $state(0)

  show(message: string, type: ToastType = 'info'): void {
    const id = this.#nextId++
    this.items = [...this.items, { id, message, type }]
    setTimeout(() => {
      this.items = this.items.filter(t => t.id !== id)
    }, 3000)
  }
}

export const toastManager = new ToastManager()
