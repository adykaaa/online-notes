import { useToast,Button } from '@chakra-ui/react'

export default function ShowToast(toast,type,msg) {
    // types are: "success", "info", "warning", "error"
        toast({
            status: type, 
            description:msg, 
            isClosable: true, 
            duration: 5000,
        })
}