import { useToast,Button } from '@chakra-ui/react'

export const CustomToast = () => {
    const toast = useToast();
    // types are: "success", "info", "warning", "error"

    const addToast = (newRes) => {
        toast({
            description:newRes.message, 
            status: newRes.type, 
            isClosable: true, 
            duration: 3500,
        })
    }
    
    return { addToast };
}