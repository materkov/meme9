import React from "react";
import * as schema from "../../schema/api";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export interface UserPageProps {
    data: schema.UserPageRenderer;
}

export class UserPage extends React.Component<UserPageProps> {
    render() {
        const {data} = this.props;

        return (
            <div>
                <Header data={data.headerRenderer}/>

                <h1>{data.name}</h1>
                <br/>
                User #{data.id}<br/>
                <Link href={data.lastPostUrl}>
                    Latest post: {data.lastPostId}
                </Link>
                <br/>
            </div>
        );
    }
}
