import React from "react";
import {Link} from "../Link/Link";
import * as schema from "../../api/api2";
import {Post} from "../Feed/Post/Post";

export function Profile(props: { data: schema.ProfileRenderer }) {
    const data = schema.ProfileRenderer.fromJSON(props.data);

    return <div>
        Profile page
        <br/>
        <Link href={"/"}>Go to feed</Link>
        <br/>
        <img alt="" src={data.avatar}/>
        <br/>
        <h2>{data.name}</h2>
        {data.posts.map(post => (
            <Post key={post.id} data={post}/>
        ))}
    </div>;
}
