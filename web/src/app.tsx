import React, { useRef, useState, useEffect } from "react";
import { VncScreen } from "react-vnc";
import "./index.css";
import axios from "axios";

interface SessionData {
	session: string;
	termination_token: string;
	created_at: string;
	started_at: string;
	valid_till: string;
}

function App() {
	// set session data in state
	const [sessionData, setSessionData] = useState<SessionData>({
		session: "",
		termination_token: "",
		created_at: "",
		started_at: "",
		valid_till: "",
	});
	const [vncUrl, setVncUrl] = useState("");
	const [terminationToken, setTerminationToken] = useState("");
	const [sessionReady, setSessionStart] = useState(false);

	const vncScreenRef = useRef<React.ElementRef<typeof VncScreen>>(null);

	const isValidSession = (vncUrl: string) => {
		if (vncUrl == "") {
			return false;
		}
		if (!vncUrl.startsWith("ws://") && !vncUrl.startsWith("wss://")) {
			return false;
		}

		return true;
	};

	// init the disposable file viewer session
	const initSqrxSession = async () => {
		try {
			// Create a disposable file viewer session
			const response = await axios.post("http://api.sqrx.com/v1/try", null);
			const data = await response.data;
			console.log(data);
			setSessionData({
				session: data.session,
				termination_token: data.termination_token,
				created_at: data.created_at,
				started_at: data.started_at,
				valid_till: data.valid_till,
			});

			console.log(sessionData);
			const wsURL = `ws://${data?.session}`;
			alert(wsURL);
			setVncUrl(`ws://${sessionData.session}`);
			setTerminationToken(data.termination_token);
			setSessionStart(true);

			return (
				<div style={{ margin: "1rem" }}>
					<VncScreen
						url={"ws://localhost:5800"}
						scaleViewport
						background="#000000"
						style={{
							width: "75vw",
							height: "75vh",
						}}
						debug
						ref={vncScreenRef}
					/>
				</div>
			);
		} catch (error) {
			console.log(error);
			alert("Error creating disposable file viewer session " + error);
		}
	};

	return (
		<>
			{/* Create Disposable File Viewer button */}
			<div
				className="d-flex justify-content-center m2p"
				data-cues="slideInDown"
				data-delay="600"
				data-disabled="true"
			>
				<span
					className="btn-span"
					data-cue="slideInDown"
					data-delay="600"
					data-show="true"
				>
					<button
						id="main-cta"
						className="btn btn-primary rounded-xl mx-1"
						onClick={initSqrxSession}
					>
						Create Disposable File Viewer
					</button>
				</span>
			</div>
			{/* End Create Disposable File Viewer button */}

			{/* connect to the VNC server via sqrx-angago reverse proxy  */}
			{/* "ws://localhost:8888/box-sqrx-rbi-1/ws" */}
			{sessionReady ? (
				<div style={{ margin: "1rem" }}>
					{isValidSession(vncUrl) ? (
						<VncScreen
							url={vncUrl}
							scaleViewport
							background="#000000"
							style={{
								width: "75vw",
								height: "75vh",
							}}
							debug
							ref={vncScreenRef}
						/>
					) : (
						<div>Sqrx Session invalid.</div>
					)}
				</div>
			) : (
				""
			)}

			{/* {isValidSession(vncUrl) ? (
				<div style={{ margin: "1rem" }}>
					<VncScreen
						url={"ws://localhost:5800"}
						scaleViewport
						background="#000000"
						style={{
							width: "75vw",
							height: "75vh",
						}}
						debug
						ref={vncScreenRef}
					/>
				</div>
			) : (
				<div>VNC URL not provided.</div>
			)} */}

			{/* destroy the session */}
			{vncScreenRef.current?.connected ? (
				<div style={{ margin: "1rem" }}>
					<button
						onClick={() => {
							const { connect, connected, disconnect } =
								vncScreenRef.current ?? {};
							if (connected) {
								disconnect?.();
								return;
							}
							connect?.();
						}}
					>
						🔥 Destroy Session
					</button>
				</div>
			) : (
				""
			)}
		</>
	);
}

export default App;
