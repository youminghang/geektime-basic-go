import { useState, useEffect } from 'react'
import axios from "@/axios/axios";

function Profile() {
    type Profile = {
        Email: string
    }
    const [data, setData] = useState<Profile>(null)
    const [isLoading, setLoading] = useState(false)

    useEffect(() => {
        setLoading(true)
        axios.get('/users/profile')
            .then((res) => res.data)
            .then((data) => {
                setData(data)
                setLoading(false)
            })
    }, [])

    if (isLoading) return <p>Loading...</p>
    if (!data) return <p>No profile data</p>

    return (
        <div>
            <h1>{data.Email}</h1>
            {/*<p>{data.bio}</p>*/}
        </div>
    )
}

export default Profile