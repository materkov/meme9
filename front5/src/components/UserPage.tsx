import React, {useEffect} from "react";
import {apiHost, Post} from "../store/types";
import {ComponentPost} from "./Post";

export function UserPage() {
    const [user, setUser] = React.useState<User>();
    const [posts, setPosts] = React.useState<Post[]>([]);
    const [loaded, setLoaded] = React.useState(false);

    useEffect(() => {
        const f = new FormData();
        f.set("id", location.pathname.substring(7));

        fetch(apiHost + "/userPage", {
            method: 'POST',
            body: f,
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                setUser(r);

                // TODO strange
                for (let post of r.posts) {
                    post.user = r;
                }

                setPosts(r.posts);
                setLoaded(true);
            })
    }, [])

    if (!loaded || !user) {
        return <>Загрузка ...</>;
    }

    return (
        <div>
            {user.name}
            <hr/>
            {posts.map(post => <ComponentPost key={post.id} post={post}/>)}
        </div>
    )
}
