import styles from "./HAREntrySections.module.sass";
import React, {useState} from "react";
import {SyntaxHighlighter} from "../SyntaxHighlighter/index";
import CollapsibleContainer from "../CollapsibleContainer";
import FancyTextDisplay from "../FancyTextDisplay";
import Checkbox from "../Checkbox";
import ProtobufDecoder from "protobuf-decoder";
var jp = require('jsonpath');

interface HAREntryViewLineProps {
    label: string;
    value: number | string;
}

const HAREntryViewLine: React.FC<HAREntryViewLineProps> = ({label, value}) => {
    return (label && value && <tr className={styles.dataLine}>
                <td className={styles.dataKey}>{label}</td>
                <td>
                    <FancyTextDisplay
                        className={styles.dataValue}
                        text={value}
                        applyTextEllipsis={false}
                        flipped={true}
                        displayIconOnMouseOver={true}
                    />
                </td>
            </tr>) || null;
}


interface HAREntrySectionCollapsibleTitleProps {
    title: string;
    isExpanded: boolean;
}

const HAREntrySectionCollapsibleTitle: React.FC<HAREntrySectionCollapsibleTitleProps> = ({title, isExpanded}) => {
    return <div className={styles.title}>
        <span className={`${styles.button} ${isExpanded ? styles.expanded : ''}`}>
            {isExpanded ? '-' : '+'}
        </span>
        <span>{title}</span>
    </div>
}

interface HAREntrySectionContainerProps {
    title: string;
}

export const HAREntrySectionContainer: React.FC<HAREntrySectionContainerProps> = ({title, children}) => {
    const [expanded, setExpanded] = useState(true);
    return <CollapsibleContainer
        className={styles.collapsibleContainer}
        isExpanded={expanded}
        onClick={() => setExpanded(!expanded)}
        title={<HAREntrySectionCollapsibleTitle title={title} isExpanded={expanded}/>}
    >
        {children}
    </CollapsibleContainer>
}

interface HAREntryBodySectionProps {
    content: any;
    encoding?: string;
    contentType?: string;
}

export const HAREntryBodySection: React.FC<HAREntryBodySectionProps> = ({
                                                                            content,
                                                                            encoding,
                                                                            contentType,
                                                                        }) => {
    const MAXIMUM_BYTES_TO_HIGHLIGHT = 10000; // The maximum of chars to highlight in body, in case the response can be megabytes
    const supportedLanguages = [['html', 'html'], ['json', 'json'], ['application/grpc', 'json']]; // [[indicator, languageToUse],...]
    const jsonLikeFormats = ['json'];
    const protobufFormats = ['application/grpc'];
    const [isWrapped, setIsWrapped] = useState(false);

    const formatTextBody = (body): string => {
        const chunk = body.slice(0, MAXIMUM_BYTES_TO_HIGHLIGHT);
        const bodyBuf = encoding === 'base64' ? atob(chunk) : chunk;

        try {
            if (jsonLikeFormats.some(format => content?.mimeType?.indexOf(format) > -1)) {
                return JSON.stringify(JSON.parse(bodyBuf), null, 2);
            } else if (protobufFormats.some(format => content?.mimeType?.indexOf(format) > -1)) {
                // Replace all non printable characters (ASCII)
                const protobufDecoder = new ProtobufDecoder(bodyBuf, true);
                return JSON.stringify(protobufDecoder.decode().toSimple(), null, 2);
            }
        } catch (error) {
            console.error(error);
        }
        return bodyBuf;
    }

    const getLanguage = (mimetype) => {
        const chunk = content.text?.slice(0, 100);
        if (chunk.indexOf('html') > 0 || chunk.indexOf('HTML') > 0) return supportedLanguages[0][1];
        const language = supportedLanguages.find(el => (mimetype + contentType).indexOf(el[0]) > -1);
        return language ? language[1] : 'default';
    }

    return <React.Fragment>
        {content && content.text?.length > 0 && <HAREntrySectionContainer title='Body'>
            <table>
                <tbody>
                    <HAREntryViewLine label={'Mime type'} value={content?.mimeType}/>
                    <HAREntryViewLine label={'Encoding'} value={encoding}/>
                </tbody>
            </table>

            <div style={{display: 'flex', alignItems: 'center', alignContent: 'center', margin: "5px 0"}} onClick={() => setIsWrapped(!isWrapped)}>
                <div style={{paddingTop: 3}}>
                    <Checkbox checked={isWrapped} onToggle={() => {}}/>
                </div>
                <span style={{marginLeft: '.5rem'}}>Wrap text</span>
            </div>

            <SyntaxHighlighter
                isWrapped={isWrapped}
                code={formatTextBody(content.text)}
                language={content?.mimeType ? getLanguage(content.mimeType) : 'default'}
            />
        </HAREntrySectionContainer>}
    </React.Fragment>
}

interface HAREntrySectionProps {
    title: string,
    arrayToIterate: any[],
}

export const HAREntryTableSection: React.FC<HAREntrySectionProps> = ({title, arrayToIterate}) => {
    return <React.Fragment>
        {
            arrayToIterate && arrayToIterate.length > 0 ?
                <HAREntrySectionContainer title={title}>
                    <table>
                        <tbody>
                            {arrayToIterate.map(({name, value}, index) => <HAREntryViewLine key={index} label={name}
                                                                                            value={value}/>)}
                        </tbody>
                    </table>
                </HAREntrySectionContainer> : <span/>
        }
    </React.Fragment>
}


interface HAREntryPolicySectionProps {
    title: string,
    response: any,
    latency?: number,
    arrayToIterate: any[],
}


interface HAREntryPolicySectionCollapsibleTitleProps {
    label: string;
    matched: string;
    isExpanded: boolean;
}

const HAREntryPolicySectionCollapsibleTitle: React.FC<HAREntryPolicySectionCollapsibleTitleProps> = ({label, matched, isExpanded}) => {
    return <div className={styles.title}>
        <span className={`${styles.button} ${isExpanded ? styles.expanded : ''}`}>
            {isExpanded ? '-' : '+'}
        </span>
        <span>
            <tr className={styles.dataLine}>
            <td className={styles.dataKey}>{label}</td>
            <td className={styles.dataKey}>{matched}</td>
            </tr>
        </span>
    </div>
}

interface HAREntryPolicySectionContainerProps {
    label: string;
    matched: string;
    children?: any;
}

export const HAREntryPolicySectionContainer: React.FC<HAREntryPolicySectionContainerProps> = ({label, matched, children}) => {
    const [expanded, setExpanded] = useState(false);
    return <CollapsibleContainer
        className={styles.collapsibleContainer}
        isExpanded={expanded}
        onClick={() => setExpanded(!expanded)}
        title={<HAREntryPolicySectionCollapsibleTitle label={label} matched={matched} isExpanded={expanded}/>}
    >
        {children}
    </CollapsibleContainer>
}

export const HAREntryTablePolicySection: React.FC<HAREntryPolicySectionProps> = ({title, response, latency, arrayToIterate}) => {
    console.log(response)
    const base64ToJson = JSON.parse(Buffer.from(response.content.text, "base64").toString());
    return <React.Fragment>
        {
            arrayToIterate && arrayToIterate.length > 0 ?
                <>
                <HAREntrySectionContainer title={title}>
                    <table>
                        <tbody>
                            {arrayToIterate.map(({rule, matched}, index) => {
                                    return (
                                        // <HAREntryViewLine key={index} label={rule.Name} value={matched}/>
                                        <HAREntryPolicySectionContainer key={index} label={rule.Name} matched={matched ? "Matched" : "Not Matched"}>
                                            {
                                                matched ? <span className={styles.dataValue}>Rule definition matched on key <b>{rule.Key}</b> with value <span className={styles.blueColor}>{rule.Value}</span></span>
                                                : <>
                                                    <span className={styles.dataValue}>Rule definition NOT matched on key <span className={styles.blueColor}>{rule.Key}</span></span>
                                                    <tr className={styles.blueColor}>Expected: {rule.Value}</tr>
                                                    <tr className={styles.latencyNotMatched}>Received: {jp.query(base64ToJson, rule.Key)}</tr>
                                                  </>
                                            }
                                            {/* <tr className={styles.dataKey}>Latency expected: {rule.Latency} ms</tr>
                                            <tr className={styles.latencyNotMatched}>Latency received: {latency} ms</tr> */}
                                            
                                        </HAREntryPolicySectionContainer>
                                    )
                                }
                            )
                            }
                        </tbody>
                    </table>
                </HAREntrySectionContainer>
                                            
                </> : <span/>
        }
    </React.Fragment>
}
