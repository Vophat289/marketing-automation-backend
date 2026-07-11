import { useState } from 'react'
import axios from 'axios'

const API_BASE = 'http://localhost:8080/api'

function App() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [token, setToken] = useState(localStorage.getItem('token') || '')

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    try {
      const res = await axios.post(`${API_BASE}/users/login`, { email, password })
      const t = res.data.token
      setToken(t)
      localStorage.setItem('token', t)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Có lỗi xảy ra')
    }
  }

  const handleLogout = async () => {
    try {
      await axios.post(`${API_BASE}/users/logout`, {}, {
        headers: { Authorization: `Bearer ${token}` }
      })
    } finally {
      setToken('')
      localStorage.removeItem('token')
    }
  }

  if (token) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div className="bg-white p-8 rounded-lg shadow-md w-96 text-center">
          <h2 className="text-2xl font-bold mb-4 text-green-600">✅ Đã đăng nhập!</h2>
          <p className="text-gray-500 text-sm mb-6 break-all">Token: {token.substring(0, 40)}...</p>
          <button
            onClick={handleLogout}
            className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-6 rounded w-full"
          >
            Đăng xuất
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-96">
        <h1 className="text-2xl font-bold mb-2 text-center text-gray-800">Marketing Automation</h1>
        <p className="text-center text-gray-400 text-sm mb-6">Đăng nhập để tiếp tục</p>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-1">Email</label>
            <input
              type="email"
              className="border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="example@email.com"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-1">Mật khẩu</label>
            <input
              type="password"
              className="border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>
          <button
            type="submit"
            className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded w-full transition-colors"
          >
            Đăng nhập
          </button>
        </form>
      </div>
    </div>
  )
}

export default App
