import { Link } from "react-router-dom";
import './login.css';
import { useState } from "react";
import { useToast,Button } from '@chakra-ui/react'
import axios from "axios";
import TextEditor from "./Hero"
import ShowToast from './Toast'

function Login() {
  
  const toast = useToast()
  const [logInSuccess,setLogInSuccess] = useState(false)
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [loginSuccess, setLoginSuccess] = useState(false);

  const login = () => {
    axios.post("http://localhost:8080/login" , { username: username, password: password })
    .then(response => {
      switch (response.status) {
        case "200":
          ShowToast(toast,"success","Login successful!")
          setLoginSuccess(true)
        case "401":
          ShowToast(toast,"error","Wrong password!")
        case "404":
          ShowToast(toast,"error","A user with this username and email has not been registered yet.")
        default:
          ShowToast(toast,"error","There is an error with the server, please try again later.")
      }
    })
    .catch(function (error) {
      ShowToast(toast,"error","There is an error with the server, please try again later.")
      console.log(error);
      });

  }
  return (
<>
    {loginSuccess ?
    (<TextEditor />) 
    :(
  <div class="login-box">
  <h2>Login to OnlineNoteZ!</h2>
  <div class="inputs">
    <div class="user-box">
      <input type="text" name="username" onChange={(e) => setUsername(e.target.value)}/>
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" onChange={(e) => setPassword(e.target.value)}/>
      <label>Password</label>
    </div>
    <a>
      <span></span>
      <span></span>
      <span></span>
      <span></span>
      <button class="submit-button" onClick={login}>Submit</button>
    </a>
  </div>
    <div class="signup">
      New here? <Link to="/register">Sign Up!</Link>
    </div>
</div>
  )};
  </>
)}

export default Login;

