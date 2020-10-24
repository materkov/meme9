import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";

export class Index extends React.Component<schema.IndexRenderer, any> {
    render() {
        return (
            <div>
                <h1>Главная страница</h1>
                {this.props.text}
                <br/>

                Лента: <Link href={"/feed"}>/feed</Link>
            </div>
        );
    }
}
