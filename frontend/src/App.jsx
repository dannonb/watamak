// import {useState} from 'react';
import logo from "./assets/images/watamak-logo.png";
import "./App.css";
// import {Greet} from "../wailsjs/go/main/App";

// function App() {
//     const [resultText, setResultText] = useState("Please enter your name below 👇");
//     const [name, setName] = useState('');
//     const updateName = (e) => setName(e.target.value);
//     const updateResultText = (result) => setResultText(result);

//     function greet() {
//         Greet(name).then(updateResultText);
//     }

//     return (
//         <div id="App">
//             <img src={logo} id="logo" alt="logo"/>
//             <div id="result" className="result">{resultText}</div>
//             <div id="input" className="input-box">
//                 <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
//                 <button className="btn" onClick={greet}>Greet</button>
//             </div>
//         </div>
//     )
// }

// export default App

import { useState, useEffect } from "react";
import {
  UploadAndWatermark,
  GetDefaultOutputPath,
  SelectFile,
} from "../wailsjs/go/main/App";

function App() {
  const [files, setFiles] = useState([]);
  const [watermarkText, setWatermarkText] = useState("My Watermark");
  const [outputPath, setOutputPath] = useState("");
  const [processedImages, setProcessedImages] = useState([]);

  useEffect(() => {
    // Set default output path (optional)
    GetDefaultOutputPath().then(setOutputPath);
  }, []);

  const selectFiles = async () => {
    try {
      const selectedFiles = await SelectFile();

      if (selectedFiles.length > 0) {
        setFiles(selectedFiles);
      }
    } catch (error) {
      console.error("Error selecting files:", error);
    }
  };

  const handleFileChange = (event) => {
    setFiles(event.target.files);
  };

  const processImages = async () => {
    if (files.length === 0) {
      alert("Please select images first.");
      return;
    }

    // const paths = Array.from(files).map((file) => file);

    try {
      const result = await UploadAndWatermark(files, watermarkText, outputPath);
      setProcessedImages(result);
    } catch (error) {
      console.error("Error processing images:", error);
    }
  };

  return (
    <div id="app">
      <div className="logo-section">
        <img src={logo} id="logo" alt="logo" />
        <h1>Watamak</h1>
      </div>

      {/* <input
        type="file"
        multiple
        accept="image/*"
        onChange={handleFileChange}
      /> */}
      <div className="image-selection">
      <button onClick={selectFiles}>Select Images</button>
      <p>
        {files.length > 0
          ? `${files.length} files selected`
          : "No files selected"}
      </p>
      </div>

      <input
        type="text"
        value={watermarkText}
        onChange={(e) => setWatermarkText(e.target.value)}
        placeholder="Enter watermark text"
        className="watermark-input"
      />
      <p>Output Folder: {outputPath || "Not set"}</p>
      <button onClick={processImages}>Apply Watermark</button>

      {processedImages.length > 0 && (
        <div>
          <h2>Processed Images</h2>
          {processedImages.map((img, index) => (
            <p key={index}>{img}</p>
          ))}
        </div>
      )}
    </div>
  );
}

export default App;
