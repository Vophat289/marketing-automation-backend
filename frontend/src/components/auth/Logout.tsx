import { useState } from 'react'
import axios from 'axios'

const API_BASE = 'http://localhost:8080/api'


function Logout() {
const handleLogout = async () => {

     const [token, setToken] = useState(localStorage.getItem('token') || '')
    
      try {
        await axios.post(`${API_BASE}/users/logout`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
      } finally {
        setToken('')
        localStorage.removeItem('token')
      }
    }

  return (
    <div>
      <button onClick={handleLogout}>Đăng Xuất</button>
    </div>
  );
}

export default Logout;
