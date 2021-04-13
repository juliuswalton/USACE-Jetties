import React, { useState, useEffect } from 'react';
import axios from 'axios';
import GeoJSONLayer from "@arcgis/core/layers/GeoJSONLayer";
import VesselMap from './VesselMap';

axios.defaults.baseURL = 'http://localhost:8080'; //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
//axios.defaults.baseURL = "https://strange-tome-305601.ue.r.appspot.com/"; //<== USE THIS LINE FOR PRODUCTION

function App() {
  const [isLoaded, setIsLoaded] = useState(false);
  const [points, setPoints] = useState([]);
  const [divisionLayer, setDivisionLayer] = useState(null);
  const [districtLayer, setDistrictLayer] = useState(null)

  useEffect(() => {
    axios.get('/api/structure')
      .then(function (response) {
        var structures = response.data; //Grab response data
        var structurePoints = [];
        //For each point in the response data create a ArcGIS Point Graphic
        for (var i = 0; i < structures.length; i++) {
          var feature = {
            geometry: {
              type: "point",
              x: structures[i].Lon,
              y: structures[i].Lat
            },
            attributes: {
              ObjectID: structures[i].ID,
              Name: structures[i].Name,
              Year: structures[i].Year,
              Type: structures[i].Type,
              Length: structures[i].Length,
              Community: structures[i].Community,
              Count: structures[i].Count
            }
          };

          //Add Point to the array of points
          structurePoints.push(feature);
        }
        setPoints(structurePoints);
        setIsLoaded(true);
      });

      // Template for the Division's popup
      const divisionTemplate = {
        title: "USACE Civil Works Division",
        content: "{DIV_SYM} : {DIVISION}",
      };
      // Layer that pulls in the geoJSON data for the USACE Divisions
      const divisionLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/4cdab8820d7a4f58aaa003be63f059ac_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: divisionTemplate,
        title: "USACE Divisions",
        listMode: false,
        visible: false,
        legendEnabled: false
      });

      const districtTemplate = {
        title: "USACE Civil Works District",
        content: "{DIV_SYM} : {District}",
      };
      const districtLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/f3e0ba4566094e74910c391eb4ecc99f_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: districtTemplate,
        title: "USACE Districts",
        listMode: false,
        visible: false,
        legendEnabled: false
      });

      setDivisionLayer(divisionLayer);
      setDistrictLayer(districtLayer);
      
  }, []);

  if(!isLoaded){
    return <div>Loading...</div>
  } else {
    return <VesselMap points={points} districtLayer={districtLayer} divisionLayer={divisionLayer}/>
  }
}

export default App;