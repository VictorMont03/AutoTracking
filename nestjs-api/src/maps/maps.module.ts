import { Module } from '@nestjs/common';
import { PlacesController } from './places/places.controller';
import { PlacesService } from './places/places.service';
import { Client as GoogleMapsClient } from '@googlemaps/google-maps-services-js';
import { DirectionsService } from './directions/directions.service';
import { DirectionsController } from './directions/directions.controller';

@Module({
  controllers: [PlacesController, DirectionsController],
  providers: [PlacesService, { provide: GoogleMapsClient, useValue: new GoogleMapsClient({}) }, DirectionsService],
  exports: [DirectionsService]
})
export class MapsModule { }
