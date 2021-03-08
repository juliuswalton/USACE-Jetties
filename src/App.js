import React, { useRef, useEffect } from "react";
//import FeatureLayer from "@arcgis/core/layers/FeatureLayer";
import ArcGISMap from "@arcgis/core/Map";
//import DictionaryRenderer from "@arcgis/core/renderers/DictionaryRenderer";
import MapView from "@arcgis/core/views/MapView";
import GeoJSONLayer from "@arcgis/core/layers/GeoJSONLayer";
import LayerList from "@arcgis/core/widgets/LayerList";
import BasemapToggle from "@arcgis/core/widgets/BasemapToggle";
import esriConfig from '@arcgis/core/config.js';
import axios from 'axios';
import Graphic from "@arcgis/core/Graphic";
import FeatureLayer from "@arcgis/core/layers/FeatureLayer";

import "./App.css";
//axios.defaults.baseURL = "http://localhost:8080"; //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
axios.defaults.baseURL = "https://strange-tome-305601.ue.r.appspot.com/"; //<== USE THIS LINE FOR PRODUCTION
function App() {

  // Required: Set this property to insure assets resolve correctly.
  esriConfig.assetsPath = './assets';

  const mapDiv = useRef(null);

  useEffect(() => {
    if (mapDiv.current) {
      /**
       * Initialize application
       */

      // Template for the Division's popup
      const template = {
        title: "USACE Civil Works Division",
        content: "{DIV_SYM} : {DIVISION}",
      };
      // Layer that pulls in the geoJSON data for the USACE Divisions
      const geoLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/4cdab8820d7a4f58aaa003be63f059ac_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: template,
        title: "USACE Divisions",
        listMode: false,
        visible: false
      });

      
      // Creates the map component
      const map = new ArcGISMap({
        basemap: "topo-vector",  // initial map styling
        layers: [geoLayer]       // array of layers that sits on top of the basemap
      });

      // Template for popup for structure points
      const structureTemplate = {
        title: "{Location}",
        content: [{
          type: "fields",
          fieldInfos: [
            {
              fieldName: "Year",
              label: "Year Constructed"
            },
            {
              fieldName: "Type",
              label: "Structure Type"
            }
          ]
        }]
      }
      // Tells the structure layer how to render the points
      const pointRenderer = {
        type: "simple",
        symbol: {
          type: "simple-marker",
          size: 10,
          color: "blue",
          outline: {
            wideth: 0.5,
            color: "white"
          }
        }
      }

      // Api call to get structure data
      var structurePoints = [];
      axios.get('/api/structure')
      .then(function (response) {
        var structures = response.data; //Grab response data
        
        //For each point in the response data create a ArcGIS Point Graphic
        for(var i = 0; i < structures.length; i++){
          var feature = {
            geometry: {
              type: "point",
              x: structures[i].Lon,
              y: structures[i].Lat
            },
            attributes: {
              ObjectID: structures[i].ID,
              Location: structures[i].Location,
              Year: structures[i].Year,
              Type: structures[i].Type
            }
          };

          //Add Point to the array of points
          structurePoints.push(feature);
        }
    
        //Create layer to show the structures
        const structureLayer = new FeatureLayer({
          title: "Structures",
          source: structurePoints, //Tell the layer where to get the data for the points
          renderer: pointRenderer,
          popupTemplate: structureTemplate,
          objectIDField: "ObjectID",
          fields: [
            {
              name: "ObjectID",
              type: "oid"
            },
            {
              name: "Location",
              type: "string"
            },
            {
              name: "Year",
              type: "integer"
            },
            {
              name: "Type",
              type: "integer"
            }
          ]
        });
    
        map.add(structureLayer); //Add the layer to the base map
      });
      // Creates the MapView - Necessary for rendering the Map Object above
      const view = new MapView({
        map: map,
        container: mapDiv.current,
        // sets the initial positioning and zoom of the map
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

      //Creates the LayerList which is used to toggle visiblity of available map layers
      view.when(function () {
        var layerList = new LayerList({
          view: view,
        });
        // adds the Layer toggle to the top-right of the screen
        view.ui.add(layerList, "top-right");
      });

      //BasemapToggle object is used to toggle the "style" of the basemap 
      var toggle = new BasemapToggle({
        view: view, // view that provides access to the map's 'topo-vector' basemap
        nextBasemap: "hybrid" // allows for toggling to the 'hybrid' basemap
      });
      // adds the toggle base map button display
      view.ui.add(toggle, "top-right");
    }
  }, []);

  return (<div className="mapDiv" ref={mapDiv}></div>

  );
}

export default App;