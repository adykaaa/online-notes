import React from 'react'
import { Card, CardHeader, CardBody, CardFooter,IconButton,Heading,Text,Divider } from '@chakra-ui/react'
import { CloseIcon, EditIcon } from '@chakra-ui/icons'

function NoteCard({id, title, text, handleDelete, handleUpdate}) {

  return (
    <Card align-self="center" background="white" maxW="350px" maxH="350px" borderRadius="md" boxShadow="box-shadow: rgba(0, 0, 0, 0.35) 0px 5px 15px;">
      <CardHeader>
        <Heading size='sm' isTruncated>{title}</Heading>
      </CardHeader>
      <Divider />
      <CardBody>
        <Text isTruncated>{text}</Text>
      </CardBody>
      <CardFooter justifyContent="space-between">
      <IconButton
          colorScheme='white'
          aria-label='Update note'
          onClick={handleUpdate}
          icon={<EditIcon alignSelf="left" color="blue"/>}/>
      <IconButton
          colorScheme='white'
          aria-label='Delete note'
          onClick={handleDelete}
          icon={<CloseIcon alignSelf="right" color="red"/>}/>
      </CardFooter>
    </Card>
  )
}

export default NoteCard