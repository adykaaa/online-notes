import { useState } from 'react'
import { Card, CardHeader, CardBody, CardFooter,IconButton,Heading,Text, Textarea } from '@chakra-ui/react'
import { CloseIcon, EditIcon, CheckIcon } from '@chakra-ui/icons'

function NoteCard({id, title, text, handleDelete, handleUpdate}) {

  const [editedText, setEditedText] = useState("")
  const [editedTitle, setEditedTitle] = useState("")
  const [isBeingEdited, setIsBeingEdited] = useState(false)

  return (
    <Card align-self="center" background="white" maxW="350px" maxH="350px" borderRadius="md" boxShadow="box-shadow: rgba(0, 0, 0, 0.35) 0px 5px 15px" border="solid #03e9f4">
      <CardHeader>
        {isBeingEdited ? <Textarea onChange={(e)=>setEditedTitle(e.target.value)}>{title}</Textarea> : <Heading size='sm' isTruncated>{title}</Heading>}
      </CardHeader>
      <CardBody>
      {isBeingEdited ? <Textarea onChange={(e)=>setEditedText(e.target.value)}>{text}</Textarea> : <Text isTruncated>{text}</Text>}
      </CardBody>
      <CardFooter justifyContent="space-between">
      <IconButton
          colorScheme='white'
          aria-label='Update note'
          onClick={()=> {
            setIsBeingEdited((state)=>!state)
            handleUpdate(id,editedTitle,editedText)}
          }
          icon={<EditIcon alignSelf="left" color="blue"/>}/>
      {isBeingEdited && <IconButton
          colorScheme='white'
          aria-label='Save changes'
          onClick={()=>{
            title = editedTitle
            text = editedText
            setIsBeingEdited(false)
          }}
          icon={<CheckIcon alignSelf="center" color="green"/>}/>}
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