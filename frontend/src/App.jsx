import { Route, Routes } from "react-router-dom";
import  Signup  from "./components/Signup.jsx";
import Login from "./components/Login.jsx"

function App() {
  return (
  <Routes>
    <Route path="/" element={<Login/>} />
    <Route path ="/register" element={<Signup />} />
  </Routes>
  )
}

export default App;
