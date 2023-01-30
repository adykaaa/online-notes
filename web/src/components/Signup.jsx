import { useState } from "react";
import axios from "axios";
import { Navigate } from "react-router-dom";
import { useToast,Button } from '@chakra-ui/react'

function SignUp() {
  
  const toast = useToast()
  const [emailReg, setEmailReg] = useState("")
  const [usernameReg, setUsernameReg] = useState("")
  const [passwordReg, setPasswordReg] = useState("")
  const [success,setSuccess] = useState(false)


  const register = () => {
    axios.post(`http://localhost:8080/register` , {email: emailReg, username: usernameReg, password: passwordReg})
    .then(response => {
      if (response.status === 201) {
        toast({
          title: 'Account created.',
          description: "Your account has been successfully created.",
          status: 'success',
          duration: 4000,
          isClosable: true,
        })
        setSuccess(true)
      }
    })

    .catch(function (error) {
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
