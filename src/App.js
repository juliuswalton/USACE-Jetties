import React, { useRef, useEffect } from "react";
//import FeatureLayer from "@arcgis/core/layers/FeatureLayer";
import ArcGISMap from "@arcgis/core/Map";
//import DictionaryRenderer from "@arcgis/core/renderers/DictionaryRenderer";
import MapView from "@arcgis/core/views/MapView";
import esriConfig from '@arcgis/core/config.js';

import "./App.css";

function App() {

  // Required: Set this property to insure assets resolve correctly.
  esriConfig.assetsPath = './assets'; 

  const mapDiv = useRef(null);

  useEffect(() => {
    if (mapDiv.current) {
      /**
       * Initialize application
       */
      const map = new ArcGISMap({
        basemap: "gray-vector",
      });

      const view = new MapView({
        map,
        container: mapDiv.current,
        extent: {
          spatialReference: {
            wkid: 102100,
          },
          xmax: -47458864.13281338,
          xmin: -54033671.557789356,
          ymax: 6678183.73641666,
          ymin: 1937864.9902844299,
        },
        zoom: 4
      });

      view.on("click", function(e){
        console.log(view.zoom);
        console.log(view.extent.toJSON());
      });
     
    }
  }, []);

  return <div className="mapDiv" ref={mapDiv}></div>;
}

export default App;