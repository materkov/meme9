import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export interface IndexProps {
    data: schema.IndexRenderer;
}

export class Index extends React.Component<IndexProps, any> {
    render() {
        return (
            <div>
                <Header data={this.props.data.headerRenderer}/>

                <h1>Главная страница</h1>
                {this.props.data.text}
                <br/>

                Лента: <Link href={this.props.data.feedUrl}>/feed</Link>
            </div>
        );
    }
}
