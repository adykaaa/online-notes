import {useContext, useState, useEffect} from 'react'
import axios from 'axios';
import { useToast, Button } from '@chakra-ui/react'
import ShowToast from './Toast'
import { UserContext } from "./UserContext";

function ViewNotes() {
    const toast = useToast()
    const [notes, setNotes] = useState([])
    const { user } = useContext(UserContext)

    useEffect(() => {
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
                setNotes(response.data)
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
    },[user])

    return (
      <>
      <ul>
      {notes.map((item) => (
        <li key={item.ID}>{item.ID}</li>
      ))}
      </ul>
      </>
  )}

export default ViewNotes