import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

const target = document.getElementById('koito-admin-mount')
if (!target) {
  throw new Error('Mount target #koito-admin-mount not found')
}

const app = mount(App, { target })

export default app
