import * as styles from "./Composer.module.css";
import React from "react";
import * as types from "../../api/api";
import {postsAdd} from "../../api/api";
import {useDiscoverPage} from "../../store/discoverPage";

export function Composer() {
    const [text, setText] = React.useState('');
    const [saving, setSaving] = React.useState(false);
    const discoverState = useDiscoverPage();

    const post = () => {
        if (!text) {
            return;
        }

        setSaving(true);
        postsAdd({text: text, pollId: pollId}).then(() => {
            setSaving(false);
            setText('');
            discoverState.refetch();
        })
    };

    const [question, setQuestion] = React.useState('');
    const [answers, setAnswers] = React.useState<string[]>(['']);

    const savePoll = () => {
        types.pollsAdd({question: question, answers: answers}).then(poll => {
            setPollId(poll.id);
            setPollActive(false);
            setQuestion('');
            setAnswers([]);
        })
    }

    const [pollActive, setPollActive] = React.useState(false);
    const [pollId, setPollId] = React.useState('');

    return <>
        <div className={styles.newPostContainer}>
                <textarea className={styles.newPost} placeholder="What's new today?" value={text}
                          onChange={e => setText(e.target.value)}/>

            {pollId && <>Poll attached</>}

            {!pollId && !pollActive && <a href="#" onClick={(e) => {
                e.preventDefault();
                setPollActive(true);
            }}>Add poll</a>}

            {pollActive && <>
                <input className={styles.pollAnswer} value={question} onChange={e => setQuestion(e.target.value)}
                       type="text" placeholder={"Your question"}/>

                {answers.map((answer, i) => (
                    <input className={styles.pollAnswer} type="text" placeholder={"Answer option..."} value={answer}
                           onChange={e => {
                               const answersCopy = structuredClone(answers);
                               answersCopy[i] = e.target.value;
                               setAnswers(answersCopy);
                           }}/>
                ))}

                <button onClick={() => setAnswers([...answers, ''])}>Add answer</button>
                <br/>
                <button onClick={savePoll}>Create poll</button>
            </>
            }
            <hr/>

            <button disabled={saving} onClick={post}>Post</button>
            <hr/>
        </div>
    </>
}