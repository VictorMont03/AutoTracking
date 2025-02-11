import { DirectionsService } from './directions.service';
import { Controller, Get, Query } from '@nestjs/common';

@Controller('directions')
export class DirectionsController {
	constructor(private readonly directionsService: DirectionsService) { }

	@Get()
	getDirections(
		@Query('originId') originId: string,
		@Query('destinationId') destinationId: string
	) {
		return this.directionsService.getDirections(originId, destinationId);
	}
}
