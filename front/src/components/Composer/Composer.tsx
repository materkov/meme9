import React from "react";
import * as schema from "../../schema/login";
import {Header} from "../Header/Header";
import {PostComposer} from "./PostComposer";
import {UnathorizedPlaceholder} from "./UnathorizedPlaceholder";

export function Composer(props: schema.ComposerRenderer) {
    return (
        <div>
            {props.headerRenderer && <Header data={props.headerRenderer}/>}

            {props.unathorizedText ?
                <UnathorizedPlaceholder data={props}/> :
                <PostComposer data={props}/>
            }
        </div>
    );
}
