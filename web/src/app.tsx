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
	const initRbixSession = async () => {
		try {
			// Create a disposable file viewer session
			const response = await axios.post("http://api.rbixlabs.com/v1/try", null);
			const data = response.data;
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
			// alert("connecting to session " + wsURL);
			alert("use password 'password' to connect to session");
			setVncUrl(wsURL);
			setTerminationToken(data.termination_token);
			setSessionStart(true);

			// hide the main content
			const maincon = document.getElementById("main-content");
			if (maincon) {
				maincon.className = "d-none";
			}
			const createBtn = document.getElementById("main-cta-container");
			if (createBtn) {
				createBtn.className = "d-none";
			}
		} catch (error) {
			console.log(error);
			alert("Error creating disposable file viewer session " + error);
		}
	};

	// destroy the disposable file viewer session
	const destroyRbixSession = async () => {
		try {
			// Create a disposable file viewer session
			console.log("destroy session of terminationToken: " + terminationToken);
			const response = await axios.post(
				`http://api.rbixlabs.com/v1/stop/${terminationToken}`,
				null
			);
			console.log(response.data);
			setSessionStart(false);

			const maincon = document.getElementById("main-content");
			if (maincon) {
				maincon.className = "d-show";
			}
			const createBtn = document.getElementById("main-cta-container");
			if (createBtn) {
				createBtn.className = "d-show";
			}
		} catch (error) {
			console.log(error);
			alert("Error destroying disposable file viewer session " + error);
		}
	};

	return (
		<>
			{/* Create Disposable File Viewer button */}
			<div
				id="main-cta-container"
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
						onClick={initRbixSession}
					>
						Create Disposable File Viewer
					</button>
				</span>
			</div>
			{/* End Create Disposable File Viewer button */}

			{/* connect to the VNC server via rbix-angago reverse proxy  */}
			{/* "ws://localhost:8888/box-rbix-rbi-1/ws" */}
			{/* "ws://localhost:5800/ws" */}
			{sessionReady ? (
				// hide the main CTA button

				<div style={{ margin: "1rem" }}>
					{isValidSession(vncUrl) ? (
						<VncScreen
							url={vncUrl}
							scaleViewport={true}
							background="#605dba"
							style={{
								height: "75vh",
							}}
							debug
							ref={vncScreenRef}
						/>
					) : (
						<div>RbiX Session invalid.</div>
					)}

					{/* destroy the session */}
					<div className="destroy-btn">
						<span
							className="btn-span"
							data-cue="slideInDown"
							data-delay="600"
							data-show="true"
						>
							<button
								className="btn btn-primary rounded-xl mx-1"
								onClick={destroyRbixSession}
							>
								🔥 Destroy Session
							</button>
						</span>
					</div>
				</div>
			) : (
				""
			)}
		</>
	);
}

export default App;
