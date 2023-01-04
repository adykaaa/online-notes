import { Route, Routes } from "react-router-dom";
import  Signup  from "./components/Signup.jsx";
import Login from "./components/Login.jsx"
import { useContext } from "react";
import { UserContext } from "./components/UserContext";

function App() {

  const user = useContext(UserContext)

  return (
  <Routes>
    <Route path="/" element={<Login/>} />
    <Route path ="/register" element={<Signup />} />
  </Routes>
  )
}

export default App;
