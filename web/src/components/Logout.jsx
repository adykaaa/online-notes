import axios from "axios";

export default function Logout() {
    axios.post("http://localhost:8080/login",username).then(response => {
        if (response.status == 200) {
            dispatch({type: 'LOGOUT', payload:username})
            localStorage.removeItem('user')
        }
    })
}