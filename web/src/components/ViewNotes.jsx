import { useContext, useState, useEffect } from 'react'
import axios from 'axios';
import { useToast, Container, SimpleGrid, Input } from '@chakra-ui/react'
import ShowToast from './Toast'
import { UserContext } from "./UserContext";
import NoteCard from './Note';
import ProSidebar from './Sidebar';

function ViewNotes() {
    const toast = useToast()
    const [notes, setNotes] = useState([])
    const { user } = useContext(UserContext)
    const [searchText,setSearchText] = useState('')


    const handleSearchTextChange = (e) => {
      setSearchText(e.target.value)
  }
    
    const handleDelete = (id) => {
      axios.delete(`http://localhost:8080/notes/${id}`,{ withCredentials: true })
      .then(response => {
          if (response.status === 200) {
            setNotes(notes.filter((note)=>note.ID !== id))
          }
          if (response.status !== 200) {
            ShowToast(toast,"error","Error deleting note, please try again!")
      }})
      .catch(function () {
        ShowToast(toast,"error","Error deleting note, please try again!")
    })}

    useEffect(() => {
        axios.get('http://localhost:8080/notes',
        {
          params:
          {
            username: user
          },
          withCredentials: true
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
      <Container minH="100%" minW='100%' display="flex" margin="0 0 0 0" padding="0 0 0 0" overflow="hidden">
        <ProSidebar/>
        <Input color="white" position="absolute" w="auto" minW="80px" marginLeft="16.5rem" focusBorderColor='white' placeholder='Search for text in a note...' marginTop="1.5rem" justifyContent="center" onChange={handleSearchTextChange}/>
        <SimpleGrid justify-content="center" align-items="center" spacing={6} margin="15" marginRight="30" marginTop="5rem" position="static" templateColumns='repeat(auto-fill, minmax(200px, 1fr))' w="70vw">
        {notes
        .filter((note)=>note.Text.String.toLowerCase().includes(searchText.toLowerCase()) || note.Title.toLowerCase().includes(searchText.toLowerCase()))
        .map((note) => (
          <NoteCard id={note.ID} title={note.Title} text={note.Text.String} handleDelete={()=>handleDelete(note.ID)} noteArray={notes} setNoteArray={setNotes}/>
        ))}
        </SimpleGrid>
      </Container>
      </>
  )}

export default ViewNotes