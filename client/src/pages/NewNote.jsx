import { useState,useContext } from 'react'
import { Card, CardHeader, CardBody, Button,Box,Text,Heading,Stack,StackDivider,Input,Textarea } from '@chakra-ui/react'
import axios from 'axios';
import { useToast} from '@chakra-ui/react'
import ShowToast from '../components/Toast'
import { UserContext } from "../components/UserContext";

function NewNote() {

    const [title,setTitle] = useState("")
    const [text,setText] = useState("")
    const toast = useToast()
    const { user } = useContext(UserContext)

    const handleTitleChange = (e) => {
        setTitle(e.target.value)
    }

    const handleTextChange = (e) => {
        setText(e.target.value)
    }

    const handleSubmit = () => {
        axios.post("http://localhost:8080/notes/create", {
            Title: title,
            User: user,
            Text: text,
        },{ withCredentials: true })
        .then(response => {
            if (response.status == 201) {
                ShowToast(toast,"success","Note successfully created!")
                setTitle("")
                setText("")
            }
        })
        .catch(function (error) {
            if (error.response) {
              switch (error.response.status) {
                case 400:
                  ShowToast(toast,"error","Duplicate or missing title!")
                  break
                case 401:
                  ShowToast(toast,"error","You are unauthorized!")
                  break
                case 403:
                  ShowToast(toast,"error","A note with this title already exists!")
                  break
                default:
                  ShowToast(toast,"error","There is an error with the server, please try again later.")
                  return
              }
            }
          })
    }


    return (
        <Card color="black" backgroundColor="linear-gradient(#141e30, #243b55)" border="solid #03e9f4" margin="auto" justifyContent="center" align-items="center">
            <CardHeader alignSelf="center" marginTop="20">
                <Heading size='lg' color="white" justifySelf="center">Create a new note</Heading>
            </CardHeader>

            <CardBody display="flex" justifyContent="flex" alignSelf="center" alignItems="center">
                <Stack size="lg" divider={<StackDivider size="10px" />} spacing='4'>
                <Box>
                    <Heading size='md' textTransform='uppercase' marginBottom="20px" color="white">
                    Title
                    </Heading>
                    <Textarea color="white" placeholder='Title of your note...' size="md" w={['40vw', '40vw', '30vw', '30vw']}  onChange={handleTitleChange}/>
                </Box>
                <Box>
                    <Heading size='md' textTransform='uppercase' marginBottom="20px" color="white">
                    Text
                    </Heading>
                    <Textarea color="white" placeholder='Text... 'size="lg" w={['40vw', '40vw', '30vw', '30vw']}  onChange={handleTextChange}/>
                </Box>
                <Button onClick={handleSubmit} colorScheme='green'>SUBMIT</Button>
                </Stack>
            </CardBody>
        </Card>
        )
    }

export default NewNote