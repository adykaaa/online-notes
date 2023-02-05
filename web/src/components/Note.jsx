import React from 'react'
import { Card, CardHeader, CardBody, CardFooter,Box,Text,Heading,Stack,StackDivider } from '@chakra-ui/react'

function Note() {
  return (
    <Card minW="full" color="White" backgroundColor="linear-gradient(#141e30, #243b55)">
        <CardHeader>
            <Heading size='lg'>Notes</Heading>
        </CardHeader>

        <CardBody>
            <Stack divider={<StackDivider colorScheme="black" size="20px" />} spacing='4'>
            <Box>
                <Heading size='xs' textTransform='uppercase'>
                Summary
                </Heading>
                <Text pt='2' fontSize='sm'>
                View a summary of all your clients over the last month.
                </Text>
            </Box>
            <Box>
                <Heading size='xs' textTransform='uppercase'>
                Overview
                </Heading>
                <Text pt='2' fontSize='sm'>
                Check out the overview of your clients.
                </Text>
            </Box>
            <Box>
                <Heading size='xs' textTransform='uppercase'>
                Analysis
                </Heading>
                <Text pt='2' fontSize='sm'>
                See a detailed analysis of all your business clients.
                </Text>
            </Box>
            </Stack>
        </CardBody>
    </Card>
    )
  }

export default Note