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
        let inputValue = e.target.value;
        setTitle(inputValue)
    }

    const handleTextChange = (e) => {
        let inputValue = e.target.value;
        setText(inputValue)
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
                  ShowToast(toast,"error","Wrongly formatted or missing Note parameter. A title is mandatory!")
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
        <Card minW="50vw" minH="50vh" color="black" backgroundColor="white" display="flex" justifyContent="center" align-items="center" marginLeft="50" marginTop="0">
            <CardHeader alignSelf="center" marginTop="20">
                <Heading size='md'>Create a new note</Heading>
            </CardHeader>

            <CardBody display="flex" justifyContent="flex" alignSelf="center">
                <Stack size="lg" divider={<StackDivider size="10px" />} spacing='4'>
                <Box>
                    <Heading size='xs' textTransform='uppercase' marginBottom="20px">
                    Title
                    </Heading>
                    <Textarea placeholder='Title of your note...' size="md" w="20vw"  onChange={handleTitleChange}/>
                </Box>
                <Box>
                    <Heading size='xs' textTransform='uppercase' marginBottom="20px">
                    Text
                    </Heading>
                    <Textarea placeholder='Text... 'size="lg" w="20vw"  onChange={handleTextChange}/>
                </Box>
                <Button onClick={handleSubmit} colorScheme='green'>SUBMIT</Button>
                </Stack>
            </CardBody>
        </Card>
        )
    }

export default NewNote