import {useState} from 'react'
import { Card, CardHeader, CardBody, CardFooter,Box,Text,Heading,Stack,StackDivider,Input,Textarea } from '@chakra-ui/react'

function NewNote() {

    const [value,setValue] = useState("")

    const handleInputChange = (e) => {
        let inputValue = e.target.value;
        setValue(inputValue)
    }

    return (
        <Card minW="100vw" color="black" backgroundColor="white" display="flex" justifyContent="center">
            <CardHeader justifyContent="center">
                <Heading size='lg'>Create a new note</Heading>
            </CardHeader>

            <CardBody w="100vw" h="100vh" display="flex" justifyContent="flex">
                <Stack divider={<StackDivider size="10px" />} spacing='4'>
                <Box>
                    <Heading size='xs' textTransform='uppercase' marginBottom="20px">
                    Title
                    </Heading>
                    <Textarea placeholder='Title of your note' size="lg" w="60vw" onChange={handleInputChange}/>
                </Box>
                <Box>
                    <Heading size='xs' textTransform='uppercase' marginBottom="20px">
                    Text
                    </Heading>
                    <Textarea placeholder='Text... 'size="lg" w="60vw" onChange={handleInputChange}/>
                </Box>
                </Stack>
            </CardBody>
        </Card>
        )
    }

export default NewNote