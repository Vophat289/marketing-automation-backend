import { useNavigate } from "react-router-dom";

function Navbar() {
  const navigate = useNavigate();
  return (
    <nav className="navbar">
      <h2>MarketFlow</h2>

      <div>
        <a href="#home">Home</a>
        <a href="#features">Features</a>
        <a href="#about">About</a>
      </div>

      <button>Đăng nhập</button>
    </nav>
  );
}

export default Navbar;
