import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export class Index extends React.Component<schema.IndexRenderer, any> {
    render() {
        return (
            <div>
                {this.props.headerRenderer && <Header data={this.props.headerRenderer}/>}

                <h1>Главная страница</h1>
                {this.props.text}
                <br/>

                Лента: <Link href={this.props.feedUrl}>/feed</Link>
            </div>
        );
    }
}
