import { PlacesService } from './places.service';
import { Controller, Get, Query } from '@nestjs/common';

@Controller('places')
export class PlacesController {
	constructor(private readonly placesService: PlacesService) { }
	@Get()
	findPlaces(@Query('text') text: string) {
		return this.placesService.findPlaces(text);
	}
}
