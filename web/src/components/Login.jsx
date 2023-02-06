import { Link, useNavigate } from "react-router-dom";
import './login.css';
import { useContext, useState } from "react";
import { useToast,Container } from '@chakra-ui/react'
import axios from "axios";
import ShowToast from './Toast'
import { UserContext } from "./UserContext";

function Login() {
  
  const toast = useToast()
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const { dispatch } = useContext(UserContext)

  const login = () => {
    axios.post("http://localhost:8080/login" , { username: username, password: password }, { withCredentials: true })
    .then(response => {
        if (response.status == 200) {
          dispatch({type: 'LOGIN', payload: username})
          localStorage.setItem('user', username)
          ShowToast(toast,"success","Login successful!")
        }
    })

    .catch(function (error) {
      if (error.response) {
        switch (error.response.status) {
          case 401:
            ShowToast(toast,"error","Wrong password!")
            break
          case 404:
            ShowToast(toast,"error","A user with this username and email has not been registered yet.")
            break
          default:
            ShowToast(toast,"error","There is an error with the server, please try again later.")
            return
        }
      }
    })
  }
  
  return (
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
  )
}

export default Login;

