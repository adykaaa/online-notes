import { useToast,Button } from '@chakra-ui/react'

export default function ShowToast(toast,type,msg) {
    // types are: "success", "info", "warning", "error"
        toast({
            description:msg, 
            status: type, 
            isClosable: true, 
            duration: 3500,
        })
    
    return toast;
}