import { useState } from "react";
import axios from "axios";
import { Navigate, useNavigate } from "react-router-dom";
import { useToast,Button } from '@chakra-ui/react'
import ShowToast from './Toast'

function SignUp() {
  
  const navigate = useNavigate()
  const toast = useToast()
  const [emailReg, setEmailReg] = useState("")
  const [usernameReg, setUsernameReg] = useState("")
  const [passwordReg, setPasswordReg] = useState("")
  const [logInSuccess,setLogInSuccess] = useState(false)


  const register = () => {
    axios.post(`http://localhost:8080/register` , {email: emailReg, username: usernameReg, password: passwordReg})
    .then(response => {
      if (response.status === 201) {
        ShowToast(toast,"success","Your account has been successfully created!")
        setLogInSuccess(true)
        navigate("/", { replace: true })
      }

      else if (response.status === 403) {
        ShowToast(toast,"error","This e-mail or username is already in use, please use another one.")
      }

      else {
        ShowToast(toast,"error","There is an error with the server, please try again later.")
      }
    })

    .catch(function (error) {
      ShowToast(toast,"error","There is an error with the server, please try again later.")
      setTimeout(navigate("/", { replace: true }), 3500)
      console.log(error);
    });

  }

  return (
    <>
  <div class="login-box">
  <h2>Register</h2>
    <div class="user-box">
      <input type="text" name="email" required="" onChange={(e) => setEmailReg(e.target.value)}/>
      <label>E-mail address</label>
    </div>
    <div class="user-box">
      <input type="text" name="username" required="" onChange={(e) => setUsernameReg(e.target.value)}/>  
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" required="" onChange={(e) => setPasswordReg(e.target.value)}/>
      <label>Password</label>
    </div>
    <a>
      <button class="register-button" onClick={register}>Register!</button>
    </a>
</div>
</>
)};

export default SignUp;
