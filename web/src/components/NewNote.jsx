import React from 'react'
import { Card, CardHeader, CardBody, CardFooter,Box,Text,Heading,Stack,StackDivider,Input } from '@chakra-ui/react'

function NewNote() {
  return (
    <Card minW="full" color="black" backgroundColor="white">
        <CardHeader>
            <Heading size='lg'>Create a new note</Heading>
        </CardHeader>

        <CardBody>
            <Stack divider={<StackDivider colorScheme="black" size="10px" />} spacing='4'>
            <Box>
                <Heading size='xs' textTransform='uppercase' marginBottom="20px">
                Title
                </Heading>
                <Input variant='outline' placeholder='Title of your note' size="md" />
            </Box>
            <Box>
                <Heading size='xs' textTransform='uppercase' marginBottom="20px">
                Text
                </Heading>
                <Input variant='outline' placeholder='Text...' size="md" />
            </Box>
            </Stack>
        </CardBody>
    </Card>
    )
  }

export default NewNote