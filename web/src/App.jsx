import { Navigate, Route, Routes } from "react-router-dom";
import  Signup  from "./components/Signup.jsx";
import Home from "./components/Home"
import Login from "./components/Login.jsx"
import { useContext, useState, useMemo} from "react";
import { UserContext } from "./components/UserContext";

function App() {

  const [user, setUser] = useState(null);
  const userValue = useMemo(() => ({ user, setUser }), [user, setUser]);
  console.log(userValue)

  return (
  <UserContext.Provider value={userValue}>
  <Routes>
    <Route path="/" element={<Login />} />
    <Route path ="/register" element={<Signup />} />
    <Route path ="/home" element={userValue ? <Home /> : <Navigate to="/" />} />
  </Routes>
  </UserContext.Provider>
  )
}

export default App;
