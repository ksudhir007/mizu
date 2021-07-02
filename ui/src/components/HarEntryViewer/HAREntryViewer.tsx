import React, {useState} from 'react';
import styles from './HAREntryViewer.module.sass';
import Tabs from "../Tabs";
import {HAREntryTableSection, HAREntryBodySection, HAREntryTablePolicySection} from "./HAREntrySections";

const MIME_TYPE_KEY = 'mimeType';

const HAREntryDisplay: React.FC<any> = ({har, entry, isCollapsed: initialIsCollapsed, isResponseMocked}) => {
    const {request, response, timings: {receive}} = entry;
    console.log(har)
    const rulesMatched = har.log.entries[0].rulesMatched
    const TABS = [
        {tab: 'request'},
        {
            tab: 'response',
            badge: <>{isResponseMocked && <span className="smallBadge virtual mock">MOCK</span>}</>
        },
        {
            tab: 'Policies Matched',
        },
    ];

    const r = request
    const [currentTab, setCurrentTab] = useState(TABS[0].tab);

    return <div className={styles.harEntry}>

        {!initialIsCollapsed && <div className={styles.body}>
            <div className={styles.bodyHeader}>
                <Tabs tabs={TABS} currentTab={currentTab} onChange={setCurrentTab} leftAligned/>
                {r?.url && <a className={styles.endpointURL} href={r.url} target='_blank' rel="noreferrer">{r.url}</a>}
            </div>
            {
                currentTab === TABS[0].tab && <React.Fragment>
                    <HAREntryTableSection title={'Headers'} arrayToIterate={r.headers}/>

                    <HAREntryTableSection title={'Cookies'} arrayToIterate={r.cookies}/>

                    {r?.postData && <HAREntryBodySection content={r.postData} encoding={r.postData.comment} contentType={r.postData[MIME_TYPE_KEY]}/>}

                    <HAREntryTableSection title={'Query'} arrayToIterate={r.queryString}/>
                </React.Fragment>
            }
            {currentTab === TABS[1].tab && <React.Fragment>
                <HAREntryTableSection title={'Headers'} arrayToIterate={response.headers}/>

                <HAREntryBodySection content={response.content} encoding={response.content?.encoding} contentType={response.content?.mimeType}/>

                <HAREntryTableSection title={'Cookies'} arrayToIterate={response.cookies}/>
            </React.Fragment>}
            {currentTab === TABS[2].tab && <React.Fragment>
                <HAREntryTablePolicySection service={har.log.entries[0].service} title={'Policy Name'} latency={receive} response={response} arrayToIterate={rulesMatched ? rulesMatched : []}/>
            </React.Fragment>}
        </div>}
    </div>;
}

interface Props {
    harObject: any;
    className?: string;
    isResponseMocked?: boolean;
    showTitle?: boolean;
}

const HAREntryViewer: React.FC<Props> = ({harObject, className, isResponseMocked, showTitle=true}) => {
    const {log: {entries}} = harObject;
    const isCollapsed = entries.length > 1;
    return <div className={`${className ? className : ''}`}>
        {Object.keys(entries).map((entry: any, index) => <HAREntryDisplay har={harObject} isCollapsed={isCollapsed} key={index} entry={entries[entry].entry} isResponseMocked={isResponseMocked} showTitle={showTitle}/>)}
    </div>
};

export default HAREntryViewer;
