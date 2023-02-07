import React from 'react'
import { Card, CardHeader, CardBody, CardFooter,Button,Heading,Text,Divider } from '@chakra-ui/react'


function NoteCard({title, text}) {
  return (
    <Card align-self="center" background="white" maxW="350px" maxH="350px">
      <CardHeader>
        <Heading size='sm' isTruncated>{title}</Heading>
      </CardHeader>
      <Divider />
    <CardBody>
      <Text isTruncated>{text}</Text>
    </CardBody>
    <CardFooter justifyContent="space-between">
      <Button alignSelf="right" color="red">Delete</Button>
      <Button alignSelf="left" color="blue">Update</Button>
    </CardFooter>
  </Card>
  )
}

export default NoteCard