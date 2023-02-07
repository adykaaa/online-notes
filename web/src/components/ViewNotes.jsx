import {useContext} from 'react'
import axios from 'axios';
import { useToast, Button} from '@chakra-ui/react'
import ShowToast from './Toast'
import { UserContext } from "./UserContext";

function getCookie(cookieName) {
    const value = "; " + document.cookie;
    const parts = value.split("; " + cookieName + "=");
    if (parts.length === 2) return parts.pop().split(";").shift();
}

function ViewNotes() {
    const toast = useToast()
    const { user } = useContext(UserContext)
    const getNotes = () => {
        axios.get('http://localhost:8080/notes',
        {
          params:
          {
            username: user
          }
          , withCredentials: true
        })
        .then(response => {
            if (response.status === 200) {
                console.log(response)
            }
        })
        .catch(function (error) {
            if (error.response) {
              switch (error.response.status) {
                case 404:
                    ShowToast(toast,"info","You don't have any notes saved yet.")
                    return
                default:
                    ShowToast(toast,"error","There is an error with the server, please try again later.")
                    return
              }
            }
        });
    }

  return (
    <Button onClick={getNotes} colorScheme='blue'>SUBMIT</Button>
  )
}

export default ViewNotes