import { useContext, useState, useEffect } from 'react'
import axios from 'axios';
import { useToast, Container, SimpleGrid } from '@chakra-ui/react'
import ShowToast from './Toast'
import { UserContext } from "./UserContext";
import NoteCard from './Note';
import ProSidebar from './Sidebar';

function ViewNotes() {
    const toast = useToast()
    const [notes, setNotes] = useState([])
    const { user } = useContext(UserContext)
    
    const handleDelete = (id) => {
      axios.delete(`http://localhost:8080/notes/${id}`,{ withCredentials: true })
      .then(response => {
          if (response.status !== 200) {
            ShowToast(toast,"error","Error deleting note, please try again!")
      }})
      .catch(function (error) {
        ShowToast(toast,"error","Error deleting note, please try again!")
    })}

    const handleUpdate = (id) => {
      console.log("Update!")
    }


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
      <Container minH="100vh" minW='100vw' display="flex" margin="0 0 0 0" padding="0 0 0 0" overflow="hidden">
      <ProSidebar/>
      <SimpleGrid justify-content="center" align-items="center" spacing={6} margin="15" marginRight="30" templateColumns='repeat(auto-fill, minmax(200px, 1fr))' w="70vw">
      {notes.map((note) => (
        <NoteCard title={note.Title} text={note.Text.String} handleDelete={()=>handleDelete(note.ID)} handleUpdate={()=>handleUpdate(note.ID)}/>
      ))}
      </SimpleGrid>
      </Container>
      </>
  )}
export default ViewNotes